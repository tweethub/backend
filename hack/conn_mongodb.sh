#!/usr/bin/env bash

if [[ -z "$mongodb_network" ]]
then
  echo "need mongodb network name"
  exit 1
fi

if [[ -z "$mongodb_user" ]]
then
  echo "need mongodb user"
  exit 2
fi

if [[ -z "$mongodb_password" ]]
then
  echo "need mongodb password"
  exit 3
fi

if [[ -z "$mongodb_image" ]]
then
  echo "need mongodb image name"
  exit 4
fi

if [[ -z "$mongodb_host" ]]
then
  echo "need mongodb host"
  exit 5
fi

MAX_RETRIES=10
retries=0

while [[ ${retries} -lt ${MAX_RETRIES} ]]
do
  echo "Trying to connect to mongodb (${retries}/${MAX_RETRIES})"

  db=$(docker run -it --rm --network "$mongodb_network" "$mongodb_image" mongo --host "$mongodb_host" -u "$mongodb_user" -p "$mongodb_password" --authenticationDatabase admin admin --quiet --eval "db.getName();" 2>/dev/null)
  if [[ "${db::-1}" == "admin" ]]; then
    exit 0
  fi

  sleep 1

  retries=$(( ${retries} + 1 ))
done

echo "failed ${retries} times to connect to mongodb"
exit 6
