# Vanilla OS SDK

The Vanilla OS SDK is a collection of libraries and tools that allow developers
to create applications for the Vanilla OS platform in a standardized and
consistent manner.

## Status

The SDK is currently in **early development**, its features are not yet stable
and are subject to change. Also note that this SDK targets Vanilla OS Vision,
not Orchid (Vanilla OS 2), so do not expect it to work in Orchid or any other
operating system.

## Run Tests

To run the tests, you can use the Containerfile provided in the repository.
This will build a container image with all the necessary dependencies and
tools to run the tests. You can build the image with the following command:


### Using the automated test script

This is the recommended way to run the tests. The script should work on both
APX and non-APX systems. It will automatically build the container image and
run the tests for you.

```bash
./test.sh
```

### Using the Containerfile

This is the manual way to run the tests. It is recommended to use this method
if for some reason the automated test script does not work for you.

#### If inside an APX container (suggested)

```bash
host-spawn podman build -t vanilla-sdk-test -f Containerfile .
```
Then, you can run the tests with the following command:

```bash
host-spawn podman run --rm --name go_vos_sdk_test_container vanilla-sdk-test
```

#### If outside an APX container

```bash
podman build -t vanilla-sdk-test -f Containerfile .
```
Then, you can run the tests with the following command:

```bash
podman run --rm --name go_vos_sdk_test_container vanilla-sdk-test
```

## Contributing

If you want to contribute to the SDK, check out the [Contributing](https://github.com/Vanilla-OS/sdk/blob/main/docs/contributing.md)
document.

## Getting Started

To get started with the SDK, check out the [Getting Started](https://github.com/Vanilla-OS/sdk/blob/main/docs/getting-started.md)
document.

## Use of Generative AI

Maintainers may use generative AI tools as assistants while working on sdk. Non-trivial assisted commits disclose the tool, model, and scope of the work.

AI tools may assist with code comments, documentation, repetitive code, and issue triage. Maintainers make project decisions and review every assisted change before it is merged.

Use these trailers for non-trivial assisted commits:

```plain
Assisted-by: <tool>:<model-version>
AI-Scope: <what the tool generated and the prompt or a short prompt summary>
```

Single-line completions, renames, and formatting changes do not need trailers.

Coding agents must also follow [AGENTS.md](AGENTS.md) before changing files,
creating commits, or opening pull requests.
