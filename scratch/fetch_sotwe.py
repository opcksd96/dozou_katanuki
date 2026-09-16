import urllib.request
import re

url = "https://www.sotwe.com/MsLuo14"
req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
try:
    with urllib.request.urlopen(req, timeout=10) as resp:
        html = resp.read().decode('utf-8', errors='ignore')
        print("Sotwe HTML len:", len(html))
        print("Carousel mentions:", len(re.findall(r'carousel', html, re.I)))
        print("Media classes:", set(re.findall(r'class="([^"]*media[^"]*)"', html, re.I)))
except Exception as e:
    print("Error fetching Sotwe directly:", e)
