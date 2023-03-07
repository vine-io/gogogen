# 安装

```shell
bash -c "$(curl -fsSL https://raw.github.com/vine-io/gogogen/master/install.sh)"
```

# deepcopy-gen
```shell
deepcopy-gen -i github.com/vine-io/apimachinery/testdata/a
```

# goproto-gen
```shell
goproto-gen --metadata-packages github.com/vine-io/apimachinery/apis/meta/v1  -p github.com/vine-io/apimachinery/testdata/a
```

# gogorm-gen
```shell
gogorm-gen  -p github.com/vine-io/apimachinery/testdata/a
```

# set-gen
```shell
 set-gen -i github.com/vine-io/gogogen/util/sets/types
```