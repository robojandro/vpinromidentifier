package vpinromidentifier_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/robojandro/vpinromidentifier"
)

func Test_VPinRomIdentifier(t *testing.T) {
	appPath, ok := os.LookupEnv("VPX_EXE_PATH")
	require.True(t, ok)

	tablesDir := "./tables/"
	tableName := "Iron Maiden (Stern 1982) v4.vpx"

	vpinRomIdentifier := vpinromidentifier.NewVPinRomIdentifier(tablesDir, appPath)
	tableRom, err := vpinRomIdentifier.ExtractTableVBS(tableName)
	assert.NoError(t, err)

	assert.Equal(t, tableName, tableRom.Table)
	assert.Equal(t, "ironmaid", tableRom.Rom)
}
