package s3_client_test

import (
	"context"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/geektheripper/vast-dsn/s3_client"
	"github.com/geektheripper/vast-dsn/s3_dsn"
)

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
		s3_opts := s3_dsn.MustNewS3("s3://minioadmin:minioadmin@localhost:9000?force-path-style=true&protocol=http")
		client, err := s3_client.NewS3Client(s3_opts)
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

		bs3_opts, bucket := s3_dsn.MustParseS3Bucket("s3://minioadmin:minioadmin@localhost:9000/foobar?force-path-style=true&protocol=http")
		bclient, err := s3_client.NewS3Client(bs3_opts)
		if err != nil {
			t.Error(err)
			return
		}
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

		os3_opts, bucket, key := s3_dsn.MustParseS3Object("s3://minioadmin:minioadmin@localhost:9000/foobar/path/to/my/object?force-path-style=true&protocol=http")
		oclient, err := s3_client.NewS3Client(os3_opts)
		if err != nil {
			t.Error(err)
			return
		}
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
