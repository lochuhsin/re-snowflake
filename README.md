<a name="readme-top"></a>


<!-- PROJECT LOGO -->
<br />
<div align="center">
  <h3 align="center">Snowflake Id Generator</h3>
  <p align="center">
    An simple, performant snowflake id generator implement with concurrency support.
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project
This project is an implementation of a Snowflake ID generator. It is extremely simple, easy to use, highly modular, and simple to modify.

Someone might wonder why I am implementing this again when there are already many available. The reason is that I just felt like it. I’m kidding. The other day, I was reading system-design related books, and it mentioned about snowflake which is sortable distributed unique identifier. So, here we are. Ha ha.

### Prerequisites
There are no prerequisites for this package.

### Installation
1. Go package manager
   ```sh
   go get -u https://github.com/lochuhsin/re-snowflake.git
   ```
2. Clone
   ```sh
    git clone https://github.com/lochuhsin/re-snowflake.git
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Performance
In general, it takes around 40ns for generating one id.
```sh
go test ./... -run=none -bench=. -benchmem -benchtime=2s -memprofile=mem.pprof -cpuprofile=cpu.pprof -blockprofile=block.pprof
goos: darwin
goarch: arm64
pkg: github.com/lochuhsin/re-snowflake
Benchmark_GenerateId-8          63802024                37.14 ns/op            0 B/op          0 allocs/op
Benchmark_GetDataCenterId-8     1000000000               0.3133 ns/op          0 B/op          0 allocs/op
Benchmark_GetMachineId-8        1000000000               0.3137 ns/op          0 B/op          0 allocs/op
Benchmark_GetSequenceNo-8       1000000000               0.3130 ns/op          0 B/op          0 allocs/op
PASS
ok      github.com/lochuhsin/re-snowflake       5.590s
```

<!-- USAGE EXAMPLES -->
## Usage

```go
    source, err := snowflake.NewSource(31, 31, 4095)
    if err != nil {
        ...
    }

    id := source.Generate()
    id.GetTime()
    id.GetDataCenterId()
    id.GetMachineId() 
    id.GetSequenceNo()
    id.GetId()

```
<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing
Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.


<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

ChuHsin(Albert) Lo - [@Linkedin](https://www.linkedin.com/in/lochuhsin/) - lochuhsin@gmail.com

Project Link: [https://github.com/lochuhsin/re-snowflake](https://github.com/lochuhsin/re-snowflake)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


