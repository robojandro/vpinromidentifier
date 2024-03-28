package vpinromidentifier

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type VPinRomIdentifier struct {
	TablesDirectory string
	VPinballAppPath string
}

// goal is to generate a JSON file mapping existing tables to the roms
// they are asking for and the reverse as well for easy look up
func NewVPinRomIdentifier(tablesDir, appPath string) *VPinRomIdentifier {
	return &VPinRomIdentifier{
		TablesDirectory: tablesDir,
		VPinballAppPath: appPath,
	}
}

type TableWithRom struct {
	Table string `json:"table"`
	Rom   string `json:"rom"`
}

func (vp *VPinRomIdentifier) ScanTables() {
}

type errorConst string

func (e errorConst) Error() string {
	return string(e)
}

const ErrRomNameNotFound errorConst = "vbs did not contain rom name"

// need to address multiple tables using the same rom
// create an alias symlink to it with
// Const cGameName = "sorcr_l2"
func (vp *VPinRomIdentifier) ExtractTableVBS(tableName string) (*TableWithRom, error) {
	woSuffix := strings.TrimSuffix(tableName, filepath.Ext(tableName))
	vbsFile := fmt.Sprintf("%s.vbs", woSuffix)

	tablePath := fmt.Sprintf("%s/%s", vp.TablesDirectory, tableName)
	if _, err := os.Stat(tablePath); err != nil {
		return nil, fmt.Errorf("did not find table file: %s\n", err)
	}

	vbsPath := fmt.Sprintf("%s/%s", vp.TablesDirectory, vbsFile)
	_, err := os.Stat(vbsPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("problem access existing vbs file: %s\n", err)
	}

	if err := vp.extractVBSFromTable(tablePath); err != nil {
		return nil, fmt.Errorf("problem extracting vbs file from vpx table '%s': %s\n", tableName, err)
	}

	fh, err := os.Open(vbsPath)
	if err != nil {
		return nil, fmt.Errorf("problem opening vbs file from vpx table '%s': %s\n", tableName, err)
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("problem scanning vbs file from vpx table '%s': %s\n", tableName, err)
	}

	var romName string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "cGameName = ") && !strings.HasPrefix(line, "'") {
			parts := strings.Split(line, "=")
			romName = strings.TrimSpace(parts[1])
			romName = strings.Trim(romName, `"`)

			break
		}
	}
	if romName == "" {
		return nil, ErrRomNameNotFound
	}
	return &TableWithRom{
		Table: tableName,
		Rom:   romName,
	}, nil
}

func (vp *VPinRomIdentifier) extractVBSFromTable(tablePath string) error {
	cmd := exec.Command(vp.VPinballAppPath, "-ExtractVBS", tablePath)
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	if strings.Contains(string(stdout), "Closing VPX...") {
		log.Printf("Finished extracting vbs for table '%s'\n", tablePath)
	}
	return nil
}
