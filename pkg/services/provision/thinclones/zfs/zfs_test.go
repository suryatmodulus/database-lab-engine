package zfs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/postgres-ai/database-lab/v2/pkg/services/provision/resources"
)

type runnerMock struct {
	cmdOutput string
	err       error
}

func (r runnerMock) Run(string, ...bool) (string, error) {
	return r.cmdOutput, r.err
}

func TestListClones(t *testing.T) {
	const (
		poolName    = "datastore"
		clonePrefix = "dblab_clone_"
	)

	testCases := []struct {
		caseName   string
		cmdOutput  string
		cloneNames []string
	}{
		{
			caseName:   "empty output",
			cloneNames: []string{},
		},
		{
			caseName: "single clone",
			cmdOutput: `datastore/clone_pre_20200831030000
datastore/dblab_clone_6000
`,
			cloneNames: []string{
				"dblab_clone_6000",
			},
		},
		{
			caseName: "multiple clones",
			cmdOutput: `datastore/clone_pre_20200831030000
datastore/dblab_clone_6000
datastore/dblab_clone_6001
`,
			cloneNames: []string{
				"dblab_clone_6000",
				"dblab_clone_6001",
			},
		},
		{
			caseName: "clone duplicate",
			cmdOutput: `datastore/clone_pre_20200831030000
datastore/dblab_clone_6000
datastore/dblab_clone_6000
`,
			cloneNames: []string{
				"dblab_clone_6000",
			},
		},
		{
			caseName: "different pool",
			cmdOutput: `datastore/clone_pre_20200831030000
dblab_pool/dblab_clone_6001
datastore/dblab_clone_6000
`,
			cloneNames: []string{
				"dblab_clone_6000",
			},
		},
		{
			caseName: "no matched clone",
			cmdOutput: `datastore/clone_pre_20200831030000
dblab_pool/dblab_clone_6001
`,
			cloneNames: []string{},
		},
	}

	for _, testCase := range testCases {
		m := Manager{
			runner: runnerMock{
				cmdOutput: testCase.cmdOutput,
			},
			config: Config{
				Pool:              resources.NewPool(poolName),
				PreSnapshotSuffix: clonePrefix,
			},
		}

		listClones, err := m.ListClonesNames()

		require.NoError(t, err, testCase.caseName)
		assert.Equal(t, testCase.cloneNames, listClones, testCase.caseName)
	}
}

func TestFailedListClones(t *testing.T) {
	m := Manager{
		runner: runnerMock{
			err: errors.New("runner error"),
		},
	}

	cloneNames, err := m.ListClonesNames()

	assert.Nil(t, cloneNames)
	assert.EqualError(t, err, "failed to list clones: runner error")
}

func TestBusySnapshotList(t *testing.T) {
	m := Manager{config: Config{Pool: &resources.Pool{Name: "dblab_pool"}}}

	out := `dblab_pool	-
dblab_pool/clone_pre_20210127105215	dblab_pool@snapshot_20210127105215_pre
dblab_pool/clone_pre_20210127113000	dblab_pool@snapshot_20210127113000_pre
dblab_pool/clone_pre_20210127120000	dblab_pool@snapshot_20210127120000_pre
dblab_pool/clone_pre_20210127123000	dblab_pool@snapshot_20210127123000_pre
dblab_pool/clone_pre_20210127130000	dblab_pool@snapshot_20210127130000_pre
dblab_pool/clone_pre_20210127133000	dblab_pool@snapshot_20210127133000_pre
dblab_pool/clone_pre_20210127140000	dblab_pool@snapshot_20210127140000_pre
dblab_pool/dblab_clone_6000	dblab_pool/clone_pre_20210127133000@snapshot_20210127133008
dblab_pool/dblab_clone_6001	dblab_pool/clone_pre_20210127123000@snapshot_20210127133008
`
	expected := []string{"dblab_pool@snapshot_20210127133000_pre", "dblab_pool@snapshot_20210127123000_pre"}

	list := m.getBusySnapshotList(out)
	require.Equal(t, 2, len(list))
	assert.Contains(t, list, expected[0])
	assert.Contains(t, list, expected[1])
}

func TestExcludingBusySnapshots(t *testing.T) {
	testCases := []struct {
		snapshotList []string
		result       string
	}{
		{
			snapshotList: []string{},
			result:       "",
		},
		{
			snapshotList: []string{"dblab_pool@snapshot_20210127133000_pre"},
			result:       "| grep -Ev 'dblab_pool@snapshot_20210127133000_pre' ",
		},
		{
			snapshotList: []string{"dblab_pool@snapshot_20210127133000_pre", "dblab_pool@snapshot_20210127123000_pre"},
			result:       "| grep -Ev 'dblab_pool@snapshot_20210127133000_pre|dblab_pool@snapshot_20210127123000_pre' ",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.result, excludeBusySnapshots(tc.snapshotList))
	}
}
