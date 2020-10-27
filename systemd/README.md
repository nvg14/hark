## Steps to install the service in ubuntu/linux

```
sudo mv hark.service /lib/systemd/system/
sudo chmod 755 /lib/systemd/system/hark.service
cd /lib/systemd/system/
sudo systemctl enable hark.service
sudo systemctl start hark
sudo journalctl -f -u hark
```