import requests
r = requests.post(
    "https://api.deepai.org/api/pose-detection",
    files={
        'image': open('squat.jpg', 'rb'),
    },
    headers={'api-key': 'quickstart-QUdJIGlzIGNvbWluZy4uLi4K'}
)
print(r.json())