# Termustat

Termustat is an online timetabling tool for university students in Iran. It works on top of [Golestan Integrated University System](https://fa.wikipedia.org/wiki/%DA%AF%D9%84%D8%B3%D8%AA%D8%A7%D9%86_(%D9%86%D8%B1%D9%85_%D8%A7%D9%81%D8%B2%D8%A7%D8%B1)), which powers major number of Iranian universities.

![](docs/screenshot.png)


## About

Termustat was developed by [Arman Jafarnezhad](https://linkedin.com/in/ArmanJ) in summer 2018 at Kharazmi University to streamline the traditional paper-based course planning process. What began as a personal project quickly evolved into a widely-adopted solution, garnering significant adoption across multiple prestigious Iranian universities including K.N. Toosi University of Technology and Shahid Beheshti University.
The platform was created to address the common challenges students face when organizing their academic schedules, offering a digital alternative to manual planning methods.

Read ["From Hesarak to Abbaspur"](https://t.me/sefroyekpub/43): The story of Termustat, published in the 7th issue of Safar-o-Yek magazine.

## How to use

This repository is an ongoing rewrite of the original Termustat. It's currently not ready for production, please check back later for a full guide. Meanwhile, you are more than welcome to contribute.

### Development

Running via Docker (to build, add `--build`)

```shell
docker compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## Architecture

### Backend

Serves as the core infrastructure, managing data models and logic while handling communication between components.

### Frontend

Provides an interactive calendar interface that enables users to
- Visualize course schedules
- Manage course selections dynamically
- View time and date overlaps

### Engine

Acts as the data processing powerhouse
- Parses course data exported from Golestan system
- Transforms raw data into structured backend models
- Ensures data compatibility and integrity

## Contribution

Pull requests are welcome.

## License

GNU General Public License v3.0 (GPL-3.0)