# revel build /home/mvilche/go/src/viltasks target -m prod
# cd target/ && ln -s src/viltasks/conf conf && \
#mkdir -p database
#cd - && 
docker build -t mvilche/viltasks .
