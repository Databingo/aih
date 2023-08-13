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
chrome_options.add_argument("--headless")
chrome_options.add_argument("user_agent=DN")
driver = uc.Chrome(options=chrome_options)

# Load cookie
driver.get("https://chat.openai.com")
with open("./4.json", "r", newline='') as inputdata:
    ck = json.load(inputdata)
for c in ck:
    driver.add_cookie({k:c[k] for k in {'name', 'value'}})

# Renew with cookie
driver.get("https://chat.openai.com")
wait = WebDriverWait(driver, 200)
try:
    notice1 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//h4[contains(text(), 'This is a free research preview.')]")))
   #print("notice1")
    next1 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//button/div[contains(text(), 'Next')]")))
   #print("next1")
    next1.click()
   #print("next1.click")
    notice2 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//h4[contains(text(), 'How we collect data')]")))
   #print("notice2")
    next2 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//button/div[contains(text(), 'Next')]")))
   #print("next2")
    next2.click()
   #print("next2.click")
    notice3 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//h4[contains(text(), 'love your feedback!')]")))
   #print("notice3")
    next3 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//button/div[contains(text(), 'Done')]")))
   #print("next3")
    next3.click()
   #print("next3.click")
    driver.find_element(By.XPATH, "//a[contains(text(), 'New chat')]").click()
    input_space = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@id='prompt-textarea']")))
   #print("login work")
except:
    print("relogin")
   #open("./2.json", "w").close()
    driver.quit()
    os.exit()

while 1:
    ori = input(":")
    if ori:
   #for line in sys.stdin:
   #    message = line.strip()
   #    ori = message.replace("(-:]", " ")
        input_space.send_keys(ori)
        driver.find_element(By.XPATH, "//button//svg:path[@d='M.5 1.163A1 1 0 0 1 1.97.28l12.868 6.837a1 1 0 0 1 0 1.766L1.969 15.72A1 1 0 0 1 .5 14.836V10.33a1 1 0 0 1 .816-.983L8.5 8 1.316 6.653A1 1 0 0 1 .5 5.67V1.163Z']").click()
       #ini_source = driver.page_source
        if ori:
            try:
                retry_icon = wait.until(EC.presence_of_element_located((By.XPATH,  "//svg:path[@d='M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15']")))
               #print("get retry_icon")
                content = retry_icon.find_element(By.XPATH,  "(//div[contains(@class, 'group w-full')])[last()]")
                text = content.get_attribute("textContent")
                text = text.replace("ChatGPTChatGPT","")
                text = text.replace("1 / 1","")
                text = text.replace("\n","(-:]")
                print(text)
                sys.stdout.flush()
                cookies = driver.get_cookies()
                with open("./4.json", "w", newline='') as outputdata:
                    json.dump(cookies, outputdata)

            except Exception as e:
                pass


