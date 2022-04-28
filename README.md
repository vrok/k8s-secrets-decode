# k8s-secrets-decode

It's a tool that can decode or encode values in Base64 in YAML documents with Kubernetes Secrets resources.

## Examples

Encode a file:

```
$ cat secret.yaml
apiVersion: v1
kind: Secret
data:
  foo: bar
  key: value
$ k8s-secrets-decode -encode -file secret.yaml
apiVersion: v1
data:
  foo: YmFy
  key: dmFsdWU=
kind: Secret  
```

Decode output from kubectl:

```
kubectl get secrets secret123 -o yaml | k8s-secrets-decode
```

## Usage

```
$ k8s-secrets-decode -help
Usage of ./k8s-secrets-decode:
  -encode
    	Encode instead of decode
  -file string
    	Read from file instead of stdin
  -output-file string
    	Write output to file instead of stdout
```