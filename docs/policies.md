[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# Policies

The CLI has commands to manage policies.

```bash
cm <verb> policies <options...>
```

## Help

```bash
cm <verb> policies -h
```

## Verbs

### Get clusters

Get the list of policies

```bash
cm get policies
```

Using `-o wide` will display additional columns about the policies, including Standards, Categories, and Controls
