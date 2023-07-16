#import undetected_chromedriver as uc
##from selenium import webdriver as uc
#import random,time,os,sys
#from selenium.webdriver.common.keys import Keys
#from selenium.webdriver.common.by import By
#from selenium.webdriver.support.ui import WebDriverWait
#from selenium.webdriver.support    import expected_conditions as EC
#import json
#import sys
#
##login = sys.argv[1]
##Login 
##if login == "login":
##    chrome_options = uc.ChromeOptions()
##    chrome_options.add_argument("--disable-extensions")
##    chrome_options.add_argument("--disable-popup-blocking")
##    chrome_options.add_argument("--profile-directory=Default")
##    chrome_options.add_argument("--ignore-certificate-errors")
##    chrome_options.add_argument("--disable-plugins-discovery")
##    chrome_options.add_argument("--incognito")
##    chrome_options.add_argument("user_agent=DN")
##    driver = uc.Chrome(options=chrome_options)
##    driver.get("https://bard.google.com")
##    #s = getpass.getpass("Press Enter after You are done login ")
##    #print("Please login google bard manually...")
##    wait = WebDriverWait(driver, 300000)
##    work = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@id='mat-input-0']")))
##    cookies = driver.get_cookies()
##    with open("./2.json", "w", newline='') as outputdata:
##        json.dump(cookies, outputdata)
##    driver.close()
#
## Restart session
##########################
##driver = uc.Chrome(options=chrome_options, headless=True)
#chrome_options = uc.ChromeOptions()
#chrome_options.add_argument("--disable-extensions")
#chrome_options.add_argument("--disable-popup-blocking")
#chrome_options.add_argument("--profile-directory=Default")
#chrome_options.add_argument("--ignore-certificate-errors")
#chrome_options.add_argument("--disable-plugins-discovery")
#chrome_options.add_argument("--incognito")
#chrome_options.add_argument("--headless")
#chrome_options.add_argument("user_agent=DN")
#driver = uc.Chrome(options=chrome_options)
#
## Load cookie
#driver.get("https://bard.google.com")
#with open("./2.json", "r", newline='') as inputdata:
#    ck = json.load(inputdata)
#for c in ck:
#    driver.add_cookie({k:c[k] for k in {'name', 'value'}})
#
## Renew with cookie
#driver.get("https://bard.google.com")
#wait = WebDriverWait(driver, 20)
#try:
#    work = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@id='mat-input-0']")))
#    print("login work")
#except:
#    print("relogin")
#    open("./2.json", "w").close()
#    driver.quit()
#    os.exit()
#
#wait = WebDriverWait(driver, 30000)
#while 1:
#   #ori = input(":")
#   #if ori:
#    for line in sys.stdin:
#        message = line.strip()
#        ori = message.replace("(-:]", " ")
#        work.send_keys(ori)
#        driver.find_element(By.XPATH, "//button[@mattooltip='Submit']").click()
#       #ini_source = driver.page_source
#        if ori:
#            try:
#                img_thinking = wait.until(EC.presence_of_element_located((By.XPATH,  "//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_thinking_v2_e272afd4f8d4bbd25efe.gif')]")))
#               #print("get img_thinking")
#                img = wait.until(EC.presence_of_element_located((By.XPATH,  "//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_resting_v2_1ff6f6a71f2d298b1a31.gif')]")))
#               #print("get img")
#                response = img.find_element(By.XPATH,  "ancestor::model-response")
#               #print("get response content img")
#                google  = response.find_element(By.XPATH,  ".//button[@aria-label='Google it']")
#                
#                contents = response.find_elements(By.XPATH, ".//message-content")
#                texts= "\n".join(content.text for content in contents)
#                text = "(-:]".join(line for line in texts.splitlines() if line)
#
#                text = response.text
#                text = text.replace("\n","(-:]")
#                text = text.replace("View other drafts","")
#                text = text.replace("Regenerate draft","")
#                text = text.replace("thumb_up","")
#                text = text.replace("thumb_down","")
#                text = text.replace("upload","")
#                text = text.replace("Google it","")
#                text = text.replace("more_vert","")
#                text = "(-:]".join(line for line in text.splitlines() if line)
#                print(text)
#                sys.stdout.flush()
#
#                cookies = driver.get_cookies()
#                with open("./2.json", "w", newline='') as outputdata:
#                    json.dump(cookies, outputdata)
#
#            except Exception as e:
#                pass
#
#
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
wait = WebDriverWait(driver, 20)
try:
    work = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@id='mat-input-0']")))
    print("login work")
except:
    print("relogin")
    open("./2.json", "w").close()
    driver.quit()
    os.exit()

wait = WebDriverWait(driver, 30000)
while 1:
   #ori = input(":")
   #if ori:
    for line in sys.stdin:
        message = line.strip()
        ori = message.replace("(-:]", " ")
        work.send_keys(ori)
        driver.find_element(By.XPATH, "//button[@mattooltip='Submit']").click()
       #ini_source = driver.page_source
        if ori:
            try:
                img_thinking = wait.until(EC.presence_of_element_located((By.XPATH,  "//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_thinking_v2_e272afd4f8d4bbd25efe.gif')]")))
               #print("get img_thinking")
                img = wait.until(EC.presence_of_element_located((By.XPATH,  "//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_resting_v2_1ff6f6a71f2d298b1a31.gif')]")))
               #print("get img")
                response = img.find_element(By.XPATH,  "ancestor::model-response")
               #print("get response content img")
                google  = response.find_element(By.XPATH,  ".//button[@aria-label='Google it']")
                
                contents = response.find_elements(By.XPATH, ".//message-content")
                texts= "\n".join(content.text for content in contents)
                text = "(-:]".join(line for line in texts.splitlines() if line)

                text = response.text
                text = text.replace("\n","(-:]")
                text = text.replace("View other drafts","")
                text = text.replace("Regenerate draft","")
                text = text.replace("thumb_up","")
                text = text.replace("thumb_down","")
                text = text.replace("upload","")
                text = text.replace("Google it","")
                text = text.replace("more_vert","")
                text = "(-:]".join(line for line in text.splitlines() if line)
                print(text)
                sys.stdout.flush()

                cookies = driver.get_cookies()
                with open("./2.json", "w", newline='') as outputdata:
                    json.dump(cookies, outputdata)

            except Exception as e:
                pass


