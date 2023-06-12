#import undetected_chromedriver as uc
#import time
#
#chrome_options = uc.ChromeOptions() ## All arguments to hide robot automation trackers
#chrome_options.add_argument("--disable-blink-features=AutomationControlled")
#chrome_options.add_argument("--no-first-run")
#chrome_options.add_argument("--no-service-autorun")
#chrome_options.add_argument("--no-default-browser-check")
#chrome_options.add_argument("--disable-extensions")
#chrome_options.add_argument("--disable-popup-blocking")
#chrome_options.add_argument("--profile-directory=Default")
#chrome_options.add_argument("--ignore-certificate-errors")
#chrome_options.add_argument("--disable-plugins-discovery")
#chrome_options.add_argument("--incognito")
#
#driver = uc.Chrome(version_main=113, options=chrome_options, headless=True)
#driver.get("https://accounts.google.com")
#driver.save_screenshot("./s.png")
#cookies = driver.get_cookies()
#print(cookies)
#time.sleep(60)
#

import undetected_chromedriver as uc
#from selenium import webdriver as uc
import random,time,os,sys
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support    import expected_conditions as EC
import json
import sys

#chrome_options.add_argument("--user-data-dir=./profile")
#driver.delete_all_cookies()

###3
#email_input.send_keys(GMAIL)
#driver.find_element(By.XPATH, "//div[@id='identifierNext']").click()
#driver.find_element(By.XPATH, "//span[text()='Next']").click()
#

login = sys.argv[1]
# Login 
if login == "login":
    chrome_options = uc.ChromeOptions()
    chrome_options.add_argument("--disable-extensions")
    chrome_options.add_argument("--disable-popup-blocking")
    chrome_options.add_argument("--profile-directory=Default")
    chrome_options.add_argument("--ignore-certificate-errors")
    chrome_options.add_argument("--disable-plugins-discovery")
    chrome_options.add_argument("--incognito")
    chrome_options.add_argument("user_agent=DN")
    driver = uc.Chrome(options=chrome_options)
    driver.get("https://bard.google.com")
    #s = getpass.getpass("Press Enter after You are done login ")
    #print("Please login google bard manually...")
    wait = WebDriverWait(driver, 300000)
    work = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@id='mat-input-0']")))
    cookies = driver.get_cookies()
    with open("./2.json", "w", newline='') as outputdata:
        json.dump(cookies, outputdata)
    driver.close()

# Restart session
#########################
#driver = uc.Chrome(options=chrome_options, headless=True)
chrome_options = uc.ChromeOptions()
chrome_options.add_argument("--disable-extensions")
chrome_options.add_argument("--disable-popup-blocking")
chrome_options.add_argument("--profile-directory=Default")
chrome_options.add_argument("--ignore-certificate-errors")
chrome_options.add_argument("--disable-plugins-discovery")
chrome_options.add_argument("--incognito")
#chrome_options.add_argument("--headless")
chrome_options.add_argument("user_agent=DN")
driver = uc.Chrome(options=chrome_options)

# Load cookie
driver.get("https://bard.google.com")
with open("./2.json", "r", newline='') as inputdata:
    ck = json.load(inputdata)
for c in ck:
    driver.add_cookie({k:c[k] for k in {'name', 'value'}})

# Renew with cookie
driver.get("https://bard.google.com")
wait = WebDriverWait(driver, 10)
try:
    work = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@id='mat-input-0']")))
    print("login work")
except:
    print("relogin clear 2.json")
    open("./2.json", "w").close()
    driver.close()
    os.exit()
    



while 1:
   #time.sleep(1)
   #print("work")
   #sys.stdout.flush()
   #line = sys.stdin.readline()
   #if not line:
   #    continue
   #message = line.strip()
   #print("Received message:", message)
#   cookies = driver.get_cookies()
#   with open("./2.json", "w", newline='') as outputdata:
#       json.dump(cookies, outputdata)

   #lines = sys.stdin.readlines()
   #print("Received message:", " ".join(lines))
    last_response_text = ""
    ori = input(":")
    if ori:
   #for line in sys.stdin:
   #    message = line.strip()
   #   #ori = message.replace("(-:]", "\n")
   #    ori = message.replace("(-:]", " ")
       #print("Received message:", message)
       #print("original message:", ori)
        work.send_keys(ori)
       #print("sended keys:", ori)
       #driver.find_element(By.XPATH, "//button[@aria-label='Send message']").click()
        driver.find_element(By.XPATH, "//button[@mattooltip='Submit']").click()
       #print("clicked")
       #st = "//user-query[text()='" + ori + "'][last()]/following-sibling::model-response/text()"
       #response = wait.until(EC.visibility_of_element_located((By.XPATH, "//model-response[last()]")))
       #e = driver.find_element(By.XPATH,  "//model-response[last()]")
       #response = wait.until(EC.visibility_of_element_located((By.XPATH,  "//message-content[last()]")))
       #content1 = driver.find_element(By.XPATH, "//message-content[last()]")
       #t1 = content1.text
        ini_source = driver.page_source
        if ori:
            try:

                img = wait.until(EC.presence_of_element_located((By.XPATH,  "//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_resting_v2_1ff6f6a71f2d298b1a31.gif')]")))
               #time.sleep(0.5)
               #time.sleep(1)
                print("get img")
               #new_source = driver.page_source
               #print("get new source")
               #while 1:
               #    if new_source != ini_source:
               #        print("page changed")
               #        break
               #response = img.find_element(By.XPATH,  "./ancestor::model-response")
               #response = img.find_element(By.XPATH,  "ancestor::model-response")
               #response  = driver.find_element(By.XPATH,  "//model-response[last()]")
               #content = response.find_element(By.XPATH, ".//message-content")
               #content = driver.find_element(By.XPATH, "//message-content[last()]")
                text = driver.find_element(By.XPATH, "//body").text
               #text = response.text 
                print(text)
               #text = "content.text"
               #while 1:
               #    content2 = driver.find_element(By.XPATH, "//message-content[last()]")
               #    if content2.text != t1:
               #        text = content2.text
               #        break

               #xp = "//user-query//h2[text()='{}']".format(ori)
               #uq = driver.find_element(By.XPATH,  xp)
                
               #contents = response.find_elements(By.XPATH, ".//message-content")
               #texts= "\n".join(content.text for content in contents)
               #text = "(-:]".join(line for line in texts.splitlines() if line)

               #text = response.text
               #text = text.replace("\n","(-:]")
               #text = text.replace("View other drafts","")
               #text = text.replace("Regenerate draft","")
               #text = text.replace("thumb_up","")
               #text = text.replace("thumb_down","")
               #text = text.replace("upload","")
               #text = text.replace("Google it","")
               #text = text.replace("more_vert","")
               #text = "(-:]".join(line for line in text.splitlines() if line)
               #with open("./bard.txt", "w") as file:
               #    file.write(text)
               #print("text[:10]")
               #print(text)
                sys.stdout.flush()

                cookies = driver.get_cookies()
                with open("./2.json", "w", newline='') as outputdata:
                    json.dump(cookies, outputdata)

            except Exception as e:
               #print(str(e))
                pass

#  #     sys.stdout.flush()
#


