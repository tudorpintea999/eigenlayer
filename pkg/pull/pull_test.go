package pull

// func TestPull(t *testing.T) {
// 	ts := []struct {
// 		name    string
// 		puller  *Puller
// 		url     string
// 		version string
// 		err     error
// 	}{
// 		{
// 			name:    "valid target",
// 			puller:  NewPuller(),
// 			url:     "https://github.com/NethermindEth/mock-avs.git",
// 			version: "v0.1.0",
// 		},
// 		{
// 			name:    "invalid version",
// 			puller:  NewPuller(),
// 			url:     "https://github.com/NethermindEth/mock-avs.git",
// 			version: "invalid-tag",
// 			err: TagNotFoundError{
// 				Tag: "invalid-tag",
// 			},
// 		},
// 		{
// 			name:    "not found or private",
// 			puller:  NewPuller(),
// 			url:     "https://github.com/NethermindEth/mock-avs-invalid.git",
// 			version: "v0.1.0",
// 			err: RepositoryNotFoundOrPrivateError{
// 				URL: "https://github.com/NethermindEth/mock-avs-invalid.git",
// 			},
// 		},
// 		// TODO: add test case using GitAuth and a private repository
// 		// TODO: add test case using GitAuth and an invalid username/password
// 		// TODO: add test case using GitAuth and an invalid url
// 	}
// 	for _, tc := range ts {
// 		t.Run(tc.name, func(t *testing.T) {
// 			_, err := tc.puller.Pull(tc.url, tc.version, t.TempDir())
// 			if tc.err == nil {
// 				assert.NoError(t, err)
// 			} else {
// 				assert.Error(t, err)
// 				assert.Equal(t, tc.err, err)
// 			}
// 		})
// 	}
// }

// func TestGetTag(t *testing.T) {
// 	testDir := t.TempDir()

// 	file, err := os.Create(filepath.Join(testDir, "test.txt"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer file.Close()

// 	if _, err := file.WriteString("Hello!"); err != nil {
// 		t.Fatal(err)
// 	}

// 	runCmdInDir(
// 		t, testDir,
// 		exec.Command("git", "init"),
// 		exec.Command("git", "config", "user.name", "user"),
// 		exec.Command("git", "config", "user.email", "user@email.com"),
// 		exec.Command("git", "add", "test.txt"),
// 		exec.Command("git", "commit", "-m", "Initial commit"),
// 		exec.Command("git", "tag", "-a", "v0.1.0", "-m", "First release"),
// 	)

// 	gitRepo, err := git.PlainOpen(testDir)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	ts := []struct {
// 		name          string
// 		tag           string
// 		expecterError error
// 	}{
// 		{
// 			name:          "valid tag",
// 			tag:           "v0.1.0",
// 			expecterError: nil,
// 		},
// 		{
// 			name:          "valid tag",
// 			tag:           "v1.0.0",
// 			expecterError: TagNotFoundError{Tag: "v1.0.0"},
// 		},
// 	}
// 	for _, tc := range ts {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tag, err := getTag(gitRepo, tc.tag)
// 			if tc.expecterError == nil {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tc.tag, tag.Name().Short())
// 			} else {
// 				assert.ErrorIs(t, tc.expecterError, err)
// 			}
// 		})
// 	}
// }

// func TestLatestTag(t *testing.T) {
// 	ts := []struct {
// 		name        string
// 		tags        []string
// 		expectedTag string
// 		err         error
// 	}{
// 		{
// 			name: "no tags",
// 			tags: []string{},
// 			err:  ErrLatestVersionNotFound,
// 		},
// 		{
// 			name: "no tag with version format",
// 			tags: []string{"v0.1", "tag1", "0.1.0"},
// 			err:  ErrLatestVersionNotFound,
// 		},
// 		{
// 			name:        "only one tag with version format",
// 			tags:        []string{"v0.1.0", "tag1", "0.1.0"},
// 			expectedTag: "v0.1.0",
// 			err:         nil,
// 		},
// 	}

// 	for _, tc := range ts {
// 		t.Run(tc.name, func(t *testing.T) {
// 			testDir := t.TempDir()

// 			file, err := os.Create(filepath.Join(testDir, "test.txt"))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer file.Close()

// 			if _, err := file.WriteString("Hello!"); err != nil {
// 				t.Fatal(err)
// 			}

// 			runCmdInDir(
// 				t, testDir,
// 				exec.Command("git", "init"),
// 				exec.Command("git", "config", "user.name", "user"),
// 				exec.Command("git", "config", "user.email", "user@email.com"),
// 				exec.Command("git", "add", "test.txt"),
// 				exec.Command("git", "commit", "-m", "Initial commit"),
// 			)

// 			var commands []*exec.Cmd
// 			for i, tag := range tc.tags {
// 				commands = append(commands, exec.Command("git", "tag", "-a", tag, "-m", fmt.Sprintf("Tag %d", i)))
// 			}
// 			runCmdInDir(t, testDir, commands...)

// 			gitRepo, err := git.PlainOpen(testDir)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			latestTag, err := getLatestVersionTag(gitRepo)
// 			if tc.err == nil {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tc.expectedTag, latestTag.Name().Short())
// 			} else {
// 				assert.Equal(t, tc.err, err)
// 				assert.Nil(t, latestTag)
// 			}
// 		})
// 	}
// }

// func runCmdInDir(t *testing.T, dir string, cmd ...*exec.Cmd) {
// 	for _, c := range cmd {
// 		t.Logf("Running command '%s' in directory %s", c, dir)
// 		c.Dir = dir
// 		if out, err := c.Output(); err != nil {
// 			t.Fatal(err)
// 		} else {
// 			t.Logf("Command output:\n%s", out)
// 		}
// 	}
// }