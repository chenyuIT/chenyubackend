echo runing chenyuAPP-contact on background.....
pkill -9 ChenyuAPP-contact
sleep 3
nohup ./ChenyuAPP-contact > ./output.log 2>&1 &

