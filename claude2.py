import undetected_chromedriver as uc
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

driver.get("https://claude.ai")

# Load cookie
with open("./3.json", "r", newline='') as inputdata:
    ck = json.load(inputdata)
for c in ck:
    driver.add_cookie({k:c[k] for k in {'name', 'value'}})

# Renew with cookie
driver.get("https://claude.ai")
wait = WebDriverWait(driver, 200)
try:
    work = wait.until(EC.visibility_of_element_located((By.XPATH,  "//p[@data-placeholder='Message Claude or search past chats...']")))
    print("login work")                                                
   #driver.find_element(By.XPATH, "//button[@class='sc-dAOort']").click()
except:
    print("relogin")
   #open("./3.json", "w").close()
    driver.quit()
    os.exit()

driver.find_element(By.XPATH, "//div[contains(text(), 'Start a new chat')]").click()
input_space = wait.until(EC.visibility_of_element_located((By.XPATH,  "//p[@data-placeholder='Message Claude...']")))

while 1:
   #ori = input(":")
   #if ori:
    for line in sys.stdin:
        message = line.strip()
        ori = message.replace("(-:]", " ")
        input_space.send_keys(ori)
        driver.find_element(By.XPATH, "//button[@aria-label='Send Message']").click()
        if ori:
            try:
                retry_icon = wait.until(EC.presence_of_element_located((By.XPATH,  "//svg:path[@d= 'M224,128a96,96,0,0,1-94.71,96H128A95.38,95.38,0,0,1,62.1,197.8a8,8,0,0,1,11-11.63A80,80,0,1,0,71.43,71.39a3.07,3.07,0,0,1-.26.25L44.59,96H72a8,8,0,0,1,0,16H24a8,8,0,0,1-8-8V56a8,8,0,0,1,16,0V85.8L60.25,60A96,96,0,0,1,224,128Z']")))
               #print("get last retry_icon")
                content = retry_icon.find_element(By.XPATH,  "preceding::div[2]")
                text = content.get_attribute("textContent")
                text = text.replace("\n","(-:]")
                print(text)
                sys.stdout.flush()
                
                # Save cookie
                cookies = driver.get_cookies()
                with open("./3.json", "w", newline='') as outputdata:
                    json.dump(cookies, outputdata)

            except Exception as e:
                pass


