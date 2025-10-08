# Titans Basketball iCal generator

Turns a fixture spreadsheet into an ical which can be imported into google calendar and then shared
with the team

## Setup for a new year

The fixture listing changes every year, but is usually a spreadsheet.

I download the "every" fixture "sheet" of the spreadsheet as TSV.

Checklist:

- Check the fields are in the same order (i.e. in main.go the `f*` constants)
- Check the Date and Time fields haven't changed format, maybe update the format string in main.go
- Run `node gen_mappings.js <input.tsv >team_mapping.go` to generate a baseline team mapping
  - This is important as we map team names to integers and so if teams change names, we can add the new name to the list for that team. This keeps the UID's deterministic even if team names change (meaning we can over-upload to google calendar and it will sort out updating matches)
- Run `go run *.go <input.tsv >output.ical` to generate the match entries.

## Calendar import

### Initial

just use the google calendar "import" function to upload the ical file.

### Updates

Deletions have to be done manually, but updates can just re-import.
So it is import to keep previous TSV's to check for matches that should be removed.
