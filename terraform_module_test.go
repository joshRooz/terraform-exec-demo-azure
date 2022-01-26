package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall" // will be removed in 0.16.0
	"github.com/rs/xid"
)

func TestTerraformResourceGroup(t *testing.T) {

	tmpDir, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir) // clean-up

	// ExactVersion also available but for CI integration tests makes
	// sense to use the latest version, if possible
	// https://pkg.go.dev/github.com/hashicorp/terraform-exec/tfinstall
	latestVersion := tfinstall.LatestVersion(tmpDir, false)
	execPath, err := tfinstall.Find(context.Background(), latestVersion)
	if err != nil {
		t.Error(err)
	}

	// Reads the configuration from ./testfixtures
	// https://pkg.go.dev/github.com/hashicorp/terraform-exec/tfexec
	tf, err := tfexec.NewTerraform("./testfixtures", execPath)
	if err != nil {
		t.Error(err)
	}

	// Log terraform executable and provider version matrix
	tfVersion, provVersions, err := tf.Version(context.Background(), false)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("terraform version: %s", tfVersion)
		for provider, version := range provVersions {
			t.Logf("%s provider version: %s", provider, version)
		}
	}

	ctx := context.Background()
	err = tf.Init(ctx, tfexec.Upgrade(false))
	if err != nil {
		t.Error(err)
	}

	// Define variables and run terraform apply
	resourceGroupName := fmt.Sprintf("resource_group_name=rg-tfexec-%s", xid.New().String())
	location := fmt.Sprint("location=australiasoutheast")

	// Ensure that terrform destroy is ran even if an error occurs
	defer tf.Destroy(ctx, tfexec.Var(resourceGroupName), tfexec.Var(location))

	t.Logf("running terraform apply, creates %s\n", resourceGroupName)
	err = tf.Apply(ctx, tfexec.Var(resourceGroupName), tfexec.Var(location))
	if err != nil {
		t.Error(err)
	}

	// more validation tests here

	// deferred terraform destroy will run now
	t.Log("running terraform destroy")
}
