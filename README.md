# macbook13-specialdeals-fetcher
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fyhinoz%2Fmacbook13-specialdeals-fetcher.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fyhinoz%2Fmacbook13-specialdeals-fetcher?ref=badge_shield)


Fetch special deals information of "macbook pro 13 inch" in json format.

## Installation

Install via `go get`:

    go get -u github.com/yhinoz/macbook13-specialdeals-fetcher

## Examples

    % macbook13-specialdeals-fetcher --locale us |jq '.[0]'
    {
      "id": "FPXQ2LL/A",
      "name": "Refurbished 13.3-inch MacBook Pro 2.3GHz dual-core Intel Core i5 with Retina display - Space Gray",
      "release": "Originally released June 2017",
      "processor": "2.3GHz dual-core Intel Core i5, Turbo Boost up to 3.6GHz, with 64MB of eDRAM",
      "memory": "8GB of 2133MHz LPDDR3 onboard memory",
      "storage": "128GB PCIe-based onboard SSD",
      "keyboard": "US",
      "price": "$1,099.00",
      "url": "https://www.apple.com/shop/product/FPXQ2LL/A/Refurbished-133-inch-MacBook-Pro-23GHz-dual-core-Intel-Core-i5-with-Retina-display-Space-Gray"
    }

If the locale option is not specified, using a default locale(jp).


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fyhinoz%2Fmacbook13-specialdeals-fetcher.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fyhinoz%2Fmacbook13-specialdeals-fetcher?ref=badge_large)