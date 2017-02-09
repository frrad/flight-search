def display_trip_option(input_option):
    out = ''
    out += '<h3> price: %s </h3>' % str(input_option['saleTotal'])

    data = input_option['slice']

    rows = []
    for datum in data:
        row = {'slice<br>duration': str(datum['duration'])}

        for segment in datum['segment']:
            row['segment<br>duration'] = str(segment['duration'])

            if 'connectionDuration' in segment:
                row['connection<br>duration'] = str(
                    segment['connectionDuration'])
            row['flight'] = str(segment['flight']['carrier']) + \
                str(segment['flight']['number'])

            for leg in segment['leg']:

                row['origin'] = str(leg['origin'])
                row['destination'] = str(leg['destination'])
                departure_dt = str(leg['departureTime'])
                row['departure<br>date'] = departure_dt[:10]
                row['departure<br>time'] = departure_dt[11:]
                arrival_dt = str(leg['arrivalTime'])
                row['arrival<br>date'] = arrival_dt[:10]
                row['arrival<br>time'] = arrival_dt[11:]
                rows.append(row)
                row = {}

    out += tableize(rows)

    # for datum in data:
    #     out += '<h4> leg </h4>'
    #     out += str(datum['segment'])

    return out


def default_sort(header):
    lol = ['slice<br>duration',
           'flight',
           'segment<br>duration',
           'connection<br>duration',
           'origin',
           'destination',
           'departure<br>date',
           'departure<br>time',
           'arrival<br>date',
           'arrival<br>time',
           ]
    if header in lol:
        return lol.index(header)
    return -1


# returns an html table given a list of dicts
def tableize(in_list, sort_by=default_sort):
    all_keys = set()
    for row in in_list:
        all_keys |= set(row.keys())

    key_list = sorted(list(all_keys), key=sort_by)

    out = '<table>\n'
    out += '<tr> '
    for header in key_list:
        out += '<th> %s </th>' % header
    out += ' </tr>\n'

    for row in in_list:
        out += '<tr> '
        for key in key_list:
            if key in row:
                out += '<td>%s</td> ' % row[key]
            else:
                out += '<td></td> '
        out += ' </tr>\n'

    out += '</table>'

    return out
