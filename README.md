## What this does ?

*   Downloads http://download.geonames.org/export/zip/US.zip
*   unzips
*   Reads US.txt
*   Produces SQL to insert into postal_code


#   Usage

*   Clone this repo
*   ```go build```
*   ```./PostalCodes > zips.sql```
*   ```psql test < zips.sql```

Big Thanks to http://www.geonames.org/, if you use this project don't forget to give credit to geonames.org.
