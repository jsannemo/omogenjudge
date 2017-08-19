struct InitArgs {
    // File descriptor to the read end of the command pipe
    int commandPipe;
    // The path where init should build its new root.
    string containerRoot;
};

// The entry point of the init process. The args pointer is to an InitArgs struct.
int Init(void* args);
