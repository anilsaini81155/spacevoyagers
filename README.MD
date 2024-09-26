## STEPS TO RUN THE APPLICATION

go mod tidy
go build
go run main.go

### SAMPLE CURLS

1)POST

  A) Terrestrial :
        curl -X POST http://localhost:8080/exoplanets \
        -H "Content-Type: application/json" \
        -d '{
        "name": "Kepler-22b",
        "description": "An Earth-like exoplanet",
        "distance": 600,
        "radius": 2.4,
        "mass": 5.972, 
        "type": "Terrestrial"
        }'
  
  B) GasGiant :

        curl -X POST http://localhost:8080/exoplanets \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Jupiter-like",
            "description": "A large gas giant",
            "distance": 1200,
            "radius": 11.2,
            "type": "GasGiant"
        }'

2) GET ALL

      curl -X GET http://localhost:8080/exoplanets

3) GET By id 

      curl -X GET http://localhost:8080/exoplanets/1

4) Update 

      curl -X PUT http://localhost:8080/exoplanets/1 \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Updated Kepler-22b",
            "description": "Updated description of Kepler-22b",
            "distance": 650,
            "radius": 2.5,
            "mass": 6.0,
            "type": "Terrestrial"
        }'

5) DELETE 

        curl -X DELETE http://localhost:8080/exoplanets/1

6) GET FUEL Estimation:

        curl -X GET http://localhost:8080/exoplanets/2/fuel?crewCapacity=100

7) LIST By Filter :


        curl -X GET http://localhost:8080/exoplanets?sort=distance

        curl -X GET http://localhost:8080/exoplanets?sort=name

        curl -X GET http://localhost:8080/exoplanets?type=GasGiant

        curl -X GET http://localhost:8080/exoplanets?min_distance=1000&max_distance=5000

        curl -X GET http://localhost:8080/exoplanets?type=Terrestrial&sort=distance


########### EXECUTING TEST CASES ############

 go clean -testcache


 go test ./handlers -v
