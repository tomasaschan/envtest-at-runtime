# THIS IS A PROOF-OF-CONCEPT

This repo exists to show an approach for using [setup-envtest] from code, instead of as a wrapper to your test.

The point is to be able to go from something like
```
KUBEBUILDER_ASSETS="$(go run sigs.k8s.io/controller-runtime/tools/setup-envtest@latest use "1.29.x!" -p path)" go test ./...
```
to just
```
go test ./...
```
with the test harness configuring the version to use and ensuring the binaries are in place.

To use this in a test (see [the full example] for details), simply
```go
func TestExample(t *testing.T) {
    binPath, err := setup.EnsureBinariesExist("1.29.x!")
    // check for errors...

    cfg, err := envtest.Environment{
        BinaryAssetsDirectory: binPath,
        // whatever else you need
    }.Start()
    // check for errors...

    c, err := client.New(cfg, client.Options{})
    // check for errors...

    // c is now a client connected to an envtest k8s api server!
}
```

[setup-envtest]: https://github.com/kubernetes-sigs/controller-runtime/blob/main/tools/setup-envtest/README.md
[the full example]: ./example/example_test.go

## Should I use this?

If you want to use this code, you should probably just copy it (or implement something similar to it) yourself; **don't import it from this repository.**

I'm going to push for this to be a supported thing in [setup-envtest] itself. We'll see how that goes...
