//Write test struct that tests various container configurations within the container_test.go file

package container

test struct {
	env            *config.Env
	config         *config.AppConfig
	cache          otter.Cache[string, any]

	sseServer      *sse.Server
	eventProcessor *event.EventProcessor
}


func (t *test) NewContainer() *container {
	return &container{
		env:            t.Env(),
		config:         t.Config(),
		cache:          t.Cache(),
		sseServer:      t.SSE(),
		eventProcessor: t.EventProcessor(),
	}
}

func (t *test) NewContainerLoadOutConfig() *container {
	cache, _ := cache.NewCache(t.Env())
	sseServer := sse.NewServer()
	eventProcessor := event.NewEventProcessor(t.Env(), t.Cache(), t.SSE())
	for _, x := range []string{"test1", "test2", "test3", "test4", "test5", "test6", "test7", "test8", "test9", "test10"} {
		eventProcessor.AddHandler(x, func(e *event.Event) {
			fmt.Println(e)

		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewContainerLoadOutConfig(tt.filePath)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRestConfig(t *testing.T) {
	tests := []struct {
		name      string


	for x := 0; x < 10; x++ {
		eventProcessor.AddHandler(fmt.Sprintf("test%d", x), func(e *event.Event) {
			fmt.Println(e)
		})
	}


	if container := t.NewContainer(); container != nil {
		return container
	}


	if r := t.NewContainer(); r != nil {
		return r
	}

	try := t.NewContainerLoadOutConfig()
	if try != nil {

		assert.NotNil(t, try.Env())
		assert.NotNil(t, try.Config())
		assert.NotNil(t, try.Cache())
		assert.NotNil(t, try.SSE())
		assert.NotNil(t, try.EventProcessor())
	}

	return nil

}

TestReadAllFilesInDir(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "example-dir-*")
	if err != nil {
		fmt.Println("Error creating temp directory:", err)
		return
	}
	defer os.RemoveAll(tempDir) // Clean up the directory and files after use

	// Create two empty files within the directory
	file1Path := filepath.Join(tempDir, "file1.txt")
	file2Path := filepath.Join(tempDir, "file2.txt")

	_, err = os.Create(file1Path)
	if err != nil {
		fmt.Println("Error creating file1:", err)
		return
	}

	_, err = os.Create(file2Path)
	if err != nil {
		fmt.Println("Error creating file2:", err)
		return
	}

	tests := []struct {
		name        string
		dirPath     string
		expectedLen int
	}{
		{
			name:        "happy path - directory exists",
			dirPath:     tempDir,
			expectedLen: 2,
		},
		{
			name:        "error path - directory does not exist",
			dirPath:     "/invalid/path",
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := readAllFilesInDir(tt.dirPath)
			assert.Equal(t, tt.expectedLen, len(files))
		})
	}
}
