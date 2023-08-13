import undetected_chromedriver as uc
import random,time,os,sys
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support    import expected_conditions as EC
import json
import sys

#########################
chrome_options = uc.ChromeOptions()
chrome_options.add_argument("--disable-extensions")
chrome_options.add_argument("--disable-popup-blocking")
chrome_options.add_argument("--profile-directory=Default")
chrome_options.add_argument("--ignore-certificate-errors")
chrome_options.add_argument("--disable-plugins-discovery")
chrome_options.add_argument("--incognito")
#chrome_options.add_argument("--headless")
#chrome_options.add_argument("user_agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/604.1 Edg/115.0.100.0")
chrome_options.add_argument("user_agent=DN")
driver = uc.Chrome(options=chrome_options)

driver.get("https://huggingface.co/chat")
with open("./5.json", "r", newline='') as inputdata:
    ck = json.load(inputdata)
for c in ck:
    driver.add_cookie({k:c[k] for k in {'name', 'value'}})

# Renew with cookie
driver.get("https://huggingface.co/chat")
wait = WebDriverWait(driver, 200)
try:
    input_space = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@enterkeyhint='send']")))
    print("login work")
except:
    print("relogin")
    driver.quit()
    os.exit()

while 1:
    ori = input(":")
    if ori:
   #for line in sys.stdin:
   #    message = line.strip()
   #    ori = message.replace("(-:]", " ")
        input_space = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@enterkeyhint='send']")))
        input_space.send_keys(ori)
        driver.find_element(By.XPATH, "//button//svg:path[@d='M27.71 4.29a1 1 0 0 0-1.05-.23l-22 8a1 1 0 0 0 0 1.87l8.59 3.43L19.59 11L21 12.41l-6.37 6.37l3.44 8.59A1 1 0 0 0 19 28a1 1 0 0 0 .92-.66l8-22a1 1 0 0 0-.21-1.05Z']").click()
        if ori:
            try:
                stop_icon = wait.until(EC.presence_of_element_located((By.XPATH,  "//svg:path[@d='M24 6H8a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2Z']")))
               #print("get stop_icon")
                wait.until(EC.staleness_of(stop_icon))
               #print("disappear stop_icon")
                img = driver.find_element(By.XPATH,  "(//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')])[last()]")
               #print("img")
                content = img.find_element(By.XPATH,  "following-sibling::div[1]")
                text = content.get_attribute("textContent")
                text = text.replace("\n","(-:]")
                print(text)
                sys.stdout.flush()
                cookies = driver.get_cookies()
                with open("./5.json", "w", newline='') as outputdata:
                    json.dump(cookies, outputdata)

            except Exception as e:
                pass


