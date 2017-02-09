import report
import json

with open('OGGNCE.out.json', 'r') as f:
    test_data = json.loads(''.join(f.readlines()))


out_data = '\n'.join([
    report.display_trip_option(option) for option in test_data['trips']['tripOption']
])

with open('/home/frederick/Downloads/test_out.html', 'w') as f:
    f.write(out_data)
