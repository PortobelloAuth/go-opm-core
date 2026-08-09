package libpostal_test

/*
http://localhost:4400/parse?address=123+MAIN+ST+APT+4+SPRINGFIELD+IL+62701

[
  {
    "label": "house_number",
    "value": "123"
  },
  {
    "label": "road",
    "value": "main st"
  },
  {
    "label": "unit",
    "value": "apt 4"
  },
  {
    "label": "city",
    "value": "springfield"
  },
  {
    "label": "state",
    "value": "il"
  },
  {
    "label": "postcode",
    "value": "62701"
  }
]
*/
