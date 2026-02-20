package mctest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func EnvTest(t testing.TB, env *envtest.Environment) *envtest.Environment {
	if env == nil {
		env = &envtest.Environment{}
	}

	env.DownloadBinaryAssets = true
	if env.BinaryAssetsDirectory == "" {
		pwd, err := os.Getwd()
		require.NoError(t, err, "failed to get current working directory")
		repoRoot, found := FindRepositoryRoot(pwd)
		require.True(t, found, "failed to find repository root")
		env.BinaryAssetsDirectory = filepath.Join(repoRoot, ".testbin")
	}

	t.Cleanup(func() {
		if err := env.Stop(); err != nil {
			t.Fatalf("failed to stop envtest: %v", err)
		}
	})

	_, err := env.Start()
	require.NoError(t, err)

	return env
}
