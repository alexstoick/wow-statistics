builddocker:
		docker build -t wow_statistics_build:latest -f ./Dockerfile-build .

ID=$(shell docker images -q wow_statistics_build)

rundocker:
		docker run -d --name wow_statistics_build -it $(ID) bash

mv_statistics:
		mkdir build
		docker cp wow_statistics_build:/go/src/github.com/alexstoick/wow-statistics/statistics ./build/statistics
build_statistics:
		docker build -t wow_statistics:latest -f ./Dockerfile-run .
tag_and_push_statistics:
		docker tag wow_statistics registry.management.stoica.xyz/wow_statistics
		docker push registry.management.stoica.xyz/wow_statistics

killdocker:
		docker stop wow_statistics_build
		docker rm wow_statistics_build

deploy_statistics: mv_statistics build_statistics tag_and_push_statistics killdocker
deploy: deploy_statistics

