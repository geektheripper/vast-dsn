package s3_dsn_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/geektheripper/vast-dsn/s3_dsn"
)

func TestS3DSNParser(t *testing.T) {
	t.Run("amazon s3", func(t *testing.T) {
		client, err := s3_dsn.NewS3("s3://access_key:secret_key@-?region=us-east-2")
		if err != nil {
			t.Error(err)
			return
		}

		options := client.Options()
		if options.Region != "us-east-2" {
			t.Error("region not parsed")
			return
		}

		cred, err := options.Credentials.Retrieve(context.Background())
		if err != nil {
			t.Error(err)
			return
		}

		if cred.AccessKeyID != "access_key" {
			t.Error("access key not parsed")
			return
		}
	})

	t.Run("self sign https minio", func(t *testing.T) {
		client, err := s3_dsn.NewS3("s3://access_key:secret_key@maggie.minio.geektr.co:9000?region=")
		if err != nil {
			t.Error(err)
			return
		}

		options := client.Options()

		if options.Region != "" {
			t.Error("region not parsed")
			return
		}

		endpoint, err := options.EndpointResolverV2.ResolveEndpoint(context.Background(), s3.EndpointParameters{
			Endpoint:       options.BaseEndpoint,
			ForcePathStyle: aws.Bool(options.UsePathStyle),
		})
		if err != nil {
			t.Error(err)
			return
		}

		str := endpoint.URI.String()
		if str != "https://maggie.minio.geektr.co:9000" {
			t.Error("endpoint not parsed")
			return
		}

		cred, err := options.Credentials.Retrieve(context.Background())
		if err != nil {
			t.Error(err)
			return
		}

		if cred.AccessKeyID != "access_key" {
			t.Error("access key not parsed")
			return
		}
	})

	t.Run("s3 bucket dsn", func(t *testing.T) {
		client, bucket, err := s3_dsn.NewS3Bucket("s3://access_key:secret_key@maggie.minio.geektr.co:9000/foobar/path/to/key?region=")
		if bucket != "" || client != nil {
			t.Error("wrong bucket dsn parsed")
			return
		}

		if fmt.Sprint(err) != fmt.Sprintf("invalid s3 bucket dsn: unexpected key: %s", "path/to/key") {
			t.Error(err)
			return
		}

		_, bucket2, _ := s3_dsn.NewS3Bucket("s3://access_key:secret_key@maggie.minio.geektr.co:9000/foobar?region=")

		if bucket2 != "foobar" {
			t.Error("bucket not parsed")
			return
		}
	})

	t.Run("s3 object dsn", func(t *testing.T) {
		_, _, key, err := s3_dsn.NewS3Object("s3://access_key:secret_key@maggie.minio.geektr.co:9000/foobar/path/to/key?region=")
		if err != nil {
			t.Error(err)
			return
		}

		if key != "path/to/key" {
			t.Error("key not parsed")
			return
		}
	})
}

func runCommand(t *testing.T, name string, args ...string) bool {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		t.Logf("%s %v failed: %s", name, args, out)
		t.Fail()
		return false
	}
	return true
}

func TestMinioCompatibility(t *testing.T) {
	t.Cleanup(func() {
		runCommand(t, "docker", "rm", "-f", "hzyvcoyzrvjkfsftadczyfejisoikwun")
		runCommand(t, "docker", "volume", "rm", "hzyvcoyzrvjkfsftadczyfejisoikwun")
	})

	if !runCommand(t, "docker", "volume", "create", "hzyvcoyzrvjkfsftadczyfejisoikwun") {
		return
	}
	if !runCommand(t, "docker", "run", "--name", "hzyvcoyzrvjkfsftadczyfejisoikwun", "-d", "--rm", "-p", "9000:9000", "-v", "hzyvcoyzrvjkfsftadczyfejisoikwun:/data", "-e", "MINIO_ROOT_USER=minioadmin", "-e", "MINIO_ROOT_PASSWORD=minioadmin", "minio/minio", "server", "/data") {
		return
	}

	for {
		_, err := http.Head("http://localhost:9000")
		if err != nil {
			t.Logf("waiting for minio to start: %s", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	t.Run("minio compatibility", func(t *testing.T) {
		client, err := s3_dsn.NewS3("s3://minioadmin:minioadmin@localhost:9000?force-path-style=true&protocol=http")
		if err != nil {
			t.Error(err)
			return
		}

		_, err = client.CreateBucket(context.Background(), &s3.CreateBucketInput{
			Bucket: aws.String("foobar"),
		})
		if err != nil {
			t.Error(err)
			return
		}

		bclient, bucket := s3_dsn.MustNewS3Bucket("s3://minioadmin:minioadmin@localhost:9000/foobar?force-path-style=true&protocol=http")
		if bucket != "foobar" {
			t.Errorf("bucket not parsed: %s", bucket)
			return
		}

		_, err = bclient.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("path/to/my/object"),
			Body:   strings.NewReader("test"),
		})
		if err != nil {
			t.Errorf("put object failed: %s", err)
			return
		}

		oclient, bucket, key := s3_dsn.MustNewS3Object("s3://minioadmin:minioadmin@localhost:9000/foobar/path/to/my/object?force-path-style=true&protocol=http")
		if bucket != "foobar" {
			t.Errorf("bucket not parsed: %s", bucket)
			return
		}

		if key != "path/to/my/object" {
			t.Errorf("key not parsed: %s", key)
			return
		}

		body, err := oclient.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			t.Errorf("get object failed: %s", err)
			return
		}

		bodyBytes, err := io.ReadAll(body.Body)
		if err != nil {
			t.Errorf("read object failed: %s", err)
			return
		}

		if string(bodyBytes) != "test" {
			t.Errorf("object content not match")
			return
		}
	})
}
