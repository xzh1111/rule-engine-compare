This repository aims to test and compare the performance of different rule engines. We will primarily focus on comparing the following rule engines:

- [hyperjumptech/grule-rule-engine](https://github.com/hyperjumptech/grule-rule-engine)

- [bilibili/gengine](https://github.com/cookedsteak/gengine)



Clone this repository to your local machine:

```
git clone https://github.com/xzh1111/rule-engine-compare.git
```

In the benchmarks/ directory, we provide code and results for performance testing. You can follow these steps to run the performance tests:


Navigate to the directory of the performance test you want to run:

```
cd rule-engine-compare/benchmarks/
```

Run the performance test:
```
go test -bench ^BenchmarkGruleExecute$ -benchtime=5s -cpu=8 -benchmem

go test -bench ^BenchmarkGenginExecute$ -benchtime=5s -cpu=8 -benchmem
```

## Contributing
Contributions to this repository are welcome! If you have any suggestions or improvements, please submit an issue or open a pull request.


## License
This repository is licensed under the MIT License.
