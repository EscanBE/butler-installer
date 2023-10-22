### Nothing, just install Butler
6. Update your configuration file template at `cmd/init.go`
7. Search pattern `// TODO Sample of template` and update with your own logic here
8. Now verify
```bash
go mod tidy
make build
./build/tbotd version
# To install binary, use `make install`
```

#### Notes
- Do not delete `// Legacy TODO xxx`, this is coding convention for TODO that keeps forever to append related code
- Always `defer libapp.TryRecoverAndExecuteExitFunctionIfRecovered` on every go-routines to release resource to prevent resources leak
- When want to exit app gracefully (eg: `os.Exit`), remember to call `libapp.ExecuteExitFunction()`
- Write the following text into README.md file of the new project
> This project follows Go Application Template version x.y