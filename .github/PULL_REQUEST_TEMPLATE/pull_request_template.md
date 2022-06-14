## Pull Request template
Please, go through these steps before you submit a PR.

1. Make sure that your PR is not a duplicate.
2. If not, then make sure that:

    a. You have done your changes in a separate branch. Branches MUST have descriptive names that start with either the `fix/` or `feature/` prefixes. Good examples are: `fix/signin-issue` or `feature/issue-templates`.

    b. You have a descriptive commit message with a short title (first line).

    c. `go test -v ./...` and `go build -v .` doesn't throw any error. If it does, fix them first and amend your commit (`git commit --amend`).

    d. You have executed `gofmt -s -w .` to format all source code.

3. **After** these steps, you're ready to open a pull request.

    a. Your pull request MUST target the `main` branch on this repository.

    b. Give a descriptive title to your PR.

    c. Provide a description of your changes.

    d. Put `closes #XXXX` or `fixes #XXXX` in your comment to auto-close the issue that your PR fixes (if such).

**PLEASE REMOVE THIS TEMPLATE BEFORE SUBMITTING**
