import requests
import datetime
import sys, os

version = "1.0"
site    = "deeplooklabs.com"
banner  = f"""

       _ _                       _    
  __ _(_) |_ ___ ___ __ _ _ _ __| |_  
 / _` | |  _(_-</ -_) _` | '_/ _| ' \ 
 \__, |_|\__/__/\___\__,_|_| \__|_||_|    {version}
 |___/                                

        {site}
"""

example = "python3 gitsearch.py \"tesla.com boto language:python\""

if len(sys.argv) < 2:
    print(banner)
    print(example)
    sys.exit()

search_term = sys.argv[1]
access_token = os.getenv('GITHUB_TOKEN')
sort_by = 'updated'
headers = {'Authorization': f'Token {access_token}'}
date_now = datetime.datetime.now()
year_now = str(date_now.year)
url = f'https://api.github.com/search/code?q={search_term}&sort={sort_by}'
response = requests.get(url, headers=headers)

def main():
    if response.status_code == 200:
        data = response.json()
        
        items = data.get('items', [])
        print(banner)
        print(f'[INF] Total found: {str(len(items))}')
        
        for item in items:
            file_name = item['name']
            file_url = item['html_url']
            
            file_response = requests.get(file_url, headers=headers)
            if file_response.status_code == 200:
                file_data = file_response.json()
                created_at = file_data['payload']['repo']['createdAt']
            else:
                created_at = 'Last update not found!'
            
            if created_at.startswith(year_now):
                print("[WRN] Recent result!")
                print(f'[{created_at}] [{file_name}] {file_url}')

            else:
                print(f'[{created_at}] [{file_name}] {file_url}')
    else:
        print('[WRN] An error occurred while making the request!')

if __name__ == '__main__':
    main()


