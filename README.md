# AddrGo

A Go-based address data importer and API server that extracts address information from OpenStreetMap (OSM) data and stores it in a MySQL database.

## Features

- **OSM Data Import**: Downloads and parses OpenStreetMap PBF files to extract address data
- **Address Deduplication**: Uses MD5 hashing to prevent duplicate address entries
- **MySQL Storage**: Stores addresses with optimized indexing for fast queries


## Prerequisites

- Go 1.25.1 or later
- MySQL database
- Internet connection for downloading OSM data

## Installation

1. Clone the repository:
```bash
git clone https://github.com/dennstack/addrgo.git
cd addrgo
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file with your database configuration:
```env
DB_USER=your_username
DB_PASS=your_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=your_database
OSM_URLS=https://download.geofabrik.de/europe/germany-latest.osm.pbf,https://example.com/another-file.osm.pbf
```

## Configuration

The application uses environment variables for configuration:

- `DB_USER`: MySQL database username
- `DB_PASS`: MySQL database password
- `DB_HOST`: MySQL database host
- `DB_PORT`: MySQL database port
- `DB_NAME`: MySQL database name
- `OSM_URLS`: Comma-separated list of OSM PBF file URLs to import

## Usage

### Running the Application

```bash
go run .
```

This will:
1. Connect to the MySQL database
2. Start importing OSM data in the background
3. Launch the HTTP API server on port 8080

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the terms specified in the LICENSE file.