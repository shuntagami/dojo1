# Requirement

- docker
- docker-compose

# How to use

## Run

```
$ cd $PROJECT_ROOT/kadai1/shuntagami/
$ make run FROM=png TO=jpg DIRNAME=sample
Successfully converted /workspace/sample/sample2/sample3/dojo4.png, to /workspace/result/sample/sample2/sample3/dojo4.jpg
Successfully converted /workspace/sample/sample2/dojo3.png, to /workspace/result/sample/sample2/dojo3.jpg
Successfully converted /workspace/sample/sample4/sample5/dojo6.PNG, to /workspace/result/sample/sample4/sample5/dojo6.jpg
Successfully converted /workspace/sample/sample4/dojo2.png, to /workspace/result/sample/sample4/dojo2.jpg
Successfully converted /workspace/sample/sample4/dojo5.png, to /workspace/result/sample/sample4/dojo5.jpg
Successfully converted /workspace/sample/dojo1.png, to /workspace/result/sample/dojo1.jpg
Successfully converted /workspace/sample/dojo2.png, to /workspace/result/sample/dojo2.jpg
```

## Test

```
$ cd $PROJECT_ROOT/kadai1/shuntagami/
$ make test
ok  	github.com/shuntagami/dojo1/kadai1/shuntagami/converter	0.730s
```
