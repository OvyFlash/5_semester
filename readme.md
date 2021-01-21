## Table of contents
1. [Types](#Types)
	1. [Post](#Post)
	2. [User](#User)
	3. [Route](#Route)
	4. [Point](#Point)
2. [Methods](#Methods)
3. [Examples](#Examples)
### Types
> ### Post
```JSON
{
	"postid": 1, //number of post in database
	"routeid": 1, //number of route in datebase
	"date": 12345, //date in unix time
	"text": "hello world", //text for post
	"route": {
        "routeid": 6,
        "userid": 4,
        "date": 1607687829,
        "points": [
            {
                "pointid": 16,
                "routeid": 6,
                "lat": "50.774",
                "lon": "25.32",
                "pointindex": 1
            },
            {
                "pointid": 17,
                "routeid": 6,
                "lat": "50.775",
                "lon": "25.33",
                "pointindex": 2
            },
            {
                "pointid": 18,
                "routeid": 6,
                "lat": "50.776",
                "lon": "25.34",
                "pointindex": 3
            }
        ],
        "start": 155556,
        "finish": 155557
    }
}
```
---
> ### User
```JSON
	"userid": 1,//user identificator; present
	"username": "ovy",//username; can be skipped
	"firstname": "Oleg",//firstname; present
	"lastname": "Oleg",//;can be skipped
	"email": "oleg@oleg.oleg",//present
	"userpassword": "userpassword",//present
	"phonenumber": 012323,//;can be skipped
	"jwt" : "eadsadsapofjpsdf123ASD",//It is not saved in database
```
---
> ### Route
```JSON
{
	"routeid": 1,//number of route saved in db; can be skipped
	"userid": 1,//user that created route; can be skipped
	"date": 12345,//unix date when this route was created; can be skipped
	"points": [//points that user passed; important
		{
			"lat": "50.774",
			"lon": "25.32"
		},
		{
			"lat": "50.775",
			"lon": "25.33",
			"pointindex": 1
		},
		{
			"lat": "50.776",
			"lon": "25.34",
		}
	],
	"start": 155556,//unix time of start; important
	"finish": 155557//unix time of finish; important
}
```
---
> ### Point
```JSON
{
	"pointid": 1, //number of point in database
	"routeid": 1, //route that point is connected with
	"lat": "30.01",//x coord
	"lon": "123.01",//y coord
	"pointindex": 1//index of point in the route
}
```
---
### Methods
Path | Description | Method | Example
--|--|--|--
/start | Returns all existing posts | GET | [EXAMPLE](#EXAMPLE-START)
/login | Returns user with JWT token | POST | [EXAMPLE](#EXAMPLE-LOGIN)
/getMe | Returns current user | GET | [EXAMPLE](#EXAMPLE-GETME)
/user/create | Returns user | POST | [EXAMPLE](#EXAMPLE-USER-CREATE)
/user/update | Updates user info | PUT | [EXAMPLE](#EXAMPLE-USER-UPDATE)
/user/delete | Deletes user from DB | DELETE | -
/user/get 	 | Gets user from DB | POST | [EXAMPLE](#EXAMPLE-USER-GET)
/route/create | Creates new route | POST | [EXAMPLE](#EXAMPLE-ROUTE-CREATE)
/route/get | Gets concrete route | POST | [EXAMPLE](#EXAMPLE-ROUTE-GET)
/route/getMy | Gets all user routes | GET | [EXAMPLE](#EXAMPLE-ROUTE-GETMY)
/post/create | Creates post | POST | [EXAMPLE](#EXAMPLE-POST-CREATE)
/post/update | Updates post | PUT | [EXAMPLE](#EXAMPLE-POST-UPDATE)
/post/delete | Deletes post | DELETE | [EXAMPLE](#EXAMPLE-POST-DELETE)
/post/get | Gets post (No auth needed) | POST | [EXAMPLE](#EXAMPLE-POST-GET)

- user/create can accept image in JPEG 
- user/update can accept image in JPEG 


### EXAMPLES
### EXAMPLE START
INPUT: nothing
OUTPUT:
```JSON
[
    {
        "postid": 2,
        "routeid": 3,
        "date": 1607623882,
        "text": "hello",
        "route": {
            "routeid": 3,
            "userid": 3,
            "date": 1607642517,
            "points": [
                {
                    "pointid": 7,
                    "routeid": 3,
                    "lat": "50.774",
                    "lon": "25.32",
                    "pointindex": 1
                },
                {
                    "pointid": 8,
                    "routeid": 3,
                    "lat": "50.775",
                    "lon": "25.33",
                    "pointindex": 2
                },
                {
                    "pointid": 9,
                    "routeid": 3,
                    "lat": "50.776",
                    "lon": "25.34",
                    "pointindex": 3
                }
            ],
            "start": 155556,
            "finish": 155557
        }
    },
]
```

### EXAMPLE LOGIN
INPUT:
```JSON
{
    "email": "3@gmail.com",
    "userpassword": "12111113"
}
```
OUTPUT:
```JSON
{
    "userid": 1,
    "firstname": "adada",
    "email": "3@gmail.com",
    "jwt": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc3MjI5ODYsImlhdCI6MTYwNzYzNjU4NiwidXNlcmlkIjoxfQ.FtSM4bK1ydgsBVHbQ2MqHig_zWffOfS7jcHFTZXfclkfQ-_d7AEUMQqZCbdaPqe5bQ7whuBQgYzuexRivTP6gaW5MLaVaO1c0YNsXF9Zd-9GCHcReAlDDtVn9Fn_AdUVHGLciqiAGrzkIZi2RFbXUBPcRsXrhO0dmW2mI5p2rbkM7ytkiXj8JAjWySLnW6WTFsm4e2ox2v7lV6wYlU4vrQeSBTPensEDuLZ_p82irFINxigfc98nlwSMrgP75a2aTuXZbRBHJwQhAKcucglQwYIjyEegM0cbkgVufKD1GHZy2-Xbrc35emHvHgeGWs7h3SqWWsVu1BCbB2ltFBzVPDVOVj0ThT5s3itxBufFN1C-lGFzuqn2VHWCo3gfNwsE242IAXsBjT-w01hzLKhn2LQsCWsQpZdY8OLOZI9xkbAlwVDq21fygC3MusYazk34IyH1apeXjZbVhUzplnZoKoBHjP0pQb6oS-IOFtHKGhe5Il1V1NR2AdA9trAjhDJm_E2SbCL0E-JAqDET878DnYPj6jEIVmMDEmnwm5ihlAk9-2oDMogZAFSTBTa0v0IzPT1AhBj4lze-MHMjg0N-A1jtsnXg-wvwVxlwW6IHR_Bm2mK68ozDry8UIinU-1TYaAOC1R1HU6th-QjrTE2xeQattNEWpztNUVgj97EqlPQ"
}
```
### EXAMPLE GETME
INPUT: -
OUTPUT: 
```JSON
{
	"userid": 1,//user identificator; present
	"username": "ovy",//username; might be skipped
	"firstname": "Oleg",//firstname; present
	"lastname": "Oleg",//;might be skipped
	"email": "oleg@oleg.oleg",//present
	"phonenumber": 012323,//;might be skipped
}
```

### EXAMPLE USER CREATE
INPUT: 
```JSON
{
    "email": "4@gmail.com",
    "userpassword": "12111113",
    "firstname": "adada"
}
```
OUTPUT:
```JSON
{
    "userid": 2,
    "firstname": "adada",
    "email": "4@gmail.com"
}
```
### EXAMPLE USER UPDATE
INPUT:
```JSON
{
    "email": "4@gmail.com",
    "userpassword": "12345678",
    "firstname": "adada"
}
```
OUTPUT:
```JSON
{
    "userid": 1,
    "firstname": "adada",
    "email": "2@gmail.com"
}
```

### EXAMPLE USER GET
INPUT:
```JSON
{
    "userid": 2
}
```
OUTPUT:
```JSON
{
    "userid": 2,
    "firstname": "adada",
    "email": "4@gmail.com"
}
```

### EXAMPLE ROUTE CREATE
INPUT:
```JSON
{
	"points": [
		{
			"pointid": 1,
			"routeid": 1,
			"lat": "50.774",
			"lon": "25.32"
		},
		{
			"pointid": 2,
			"routeid": 1,
			"lat": "50.775",
			"lon": "25.33",
			"pointindex": 1
		},
		{
			"pointid": 3,
			"routeid": 1,
			"lat": "50.776",
			"lon": "25.34",
			"pointindex": 2
		}
	],
	"start": 155556,
	"finish": 155557
}

```
OUTPUT: 
```JSON
{
    "routeid": 3,
    "userid": 3,
    "date": 1607642517,
    "start": 155556,
    "finish": 155557
}
```

### EXAMPLE ROUTE GET
INPUT:
```JSON
{
    "routeid": 1
}
```
OUTPUT:
```JSON
{
    "routeid": 1,
    "userid": 1,
    "date": 1607616077,
    "points": [
        {
            "pointid": 1,
            "routeid": 1,
            "lat": "50.774",
            "lon": "25.32"
        },
        {
            "pointid": 2,
            "routeid": 1,
            "lat": "50.775",
            "lon": "25.33",
            "pointindex": 1
        },
        {
            "pointid": 3,
            "routeid": 1,
            "lat": "50.776",
            "lon": "25.34",
            "pointindex": 2
        }
    ],
    "start": 155556,
    "finish": 155557
}
```

### EXAMPLE ROUTE GETMY
INPUT: -
OUTPUT:
```JSON
[
    {
        "routeid": 3,
        "userid": 3,
        "date": 1607642517,
        "points": [
            {
                "pointid": 7,
                "routeid": 3,
                "lat": "50.774",
                "lon": "25.32",
                "pointindex": 1
            },
            {
                "pointid": 8,
                "routeid": 3,
                "lat": "50.775",
                "lon": "25.33",
                "pointindex": 2
            },
            {
                "pointid": 9,
                "routeid": 3,
                "lat": "50.776",
                "lon": "25.34",
                "pointindex": 3
            }
        ],
        "start": 155556,
        "finish": 155557
    },
    {
        "routeid": 4,
        "userid": 3,
        "date": 1607642675,
        "points": [
            {
                "pointid": 10,
                "routeid": 4,
                "lat": "50.774",
                "lon": "25.32",
                "pointindex": 1
            },
            {
                "pointid": 11,
                "routeid": 4,
                "lat": "50.775",
                "lon": "25.33",
                "pointindex": 2
            },
            {
                "pointid": 12,
                "routeid": 4,
                "lat": "50.776",
                "lon": "25.34",
                "pointindex": 3
            }
        ],
        "start": 155556,
        "finish": 155557
    }
]
```

### EXAMPLE POST CREATE
INPUT: 
```JSON
{
    "routeid": 2,
    "text": "hello"
}
```
OUTPUT:
```JSON
{
    "postid": 10,
    "routeid": 6,
    "date": 1607688138,
    "text": "hello",
    "route": {
        "routeid": 6,
        "userid": 4,
        "date": 1607687829,
        "points": [
            {
                "pointid": 16,
                "routeid": 6,
                "lat": "50.774",
                "lon": "25.32",
                "pointindex": 1
            },
            {
                "pointid": 17,
                "routeid": 6,
                "lat": "50.775",
                "lon": "25.33",
                "pointindex": 2
            },
            {
                "pointid": 18,
                "routeid": 6,
                "lat": "50.776",
                "lon": "25.34",
                "pointindex": 3
            }
        ],
        "start": 155556,
        "finish": 155557
    }
}
```

### EXAMPLE POST UPDATE
INPUT:
```JSON
{
    "postid": 8,
    "text": "goodbye"
}
```
OUTPUT:
```JSON
{
    "postid": 8,
    "routeid": 2,
    "date": 1607684797,
    "text": "goodbye",
    "route": {
        "routeid": 2,
        "userid": 1,
        "date": 1607616080,
        "points": [
            {
                "pointid": 4,
                "routeid": 2,
                "lat": "50.774",
                "lon": "25.32"
            },
            {
                "pointid": 5,
                "routeid": 2,
                "lat": "50.775",
                "lon": "25.33",
                "pointindex": 1
            },
            {
                "pointid": 6,
                "routeid": 2,
                "lat": "50.776",
                "lon": "25.34",
                "pointindex": 2
            }
        ],
        "start": 155556,
        "finish": 155557
    }
}
```

### EXAMPLE POST DELETE
INPUT:
```JSON
{
    "postid": 1
}
```
OUTPUT: -

### EXAMPLE POST GET
INPUT:
```JSON
{
    "postid": 3
}
```
OUTPUT:
```JSON
{
    "postid": 3,
    "routeid": 1,
    "date": 1607625134,
    "text": "hello",
    "route": {
        "routeid": 1,
        "userid": 1,
        "date": 1607616077,
        "points": [
            {
                "pointid": 1,
                "routeid": 1,
                "lat": "50.774",
                "lon": "25.32"
            },
            {
                "pointid": 2,
                "routeid": 1,
                "lat": "50.775",
                "lon": "25.33",
                "pointindex": 1
            },
            {
                "pointid": 3,
                "routeid": 1,
                "lat": "50.776",
                "lon": "25.34",
                "pointindex": 2
            }
        ],
        "start": 155556,
        "finish": 155557
    }
}
```




- /follow подписаться "POST"
- /unfollow отписаться. Для них следующий "DELETE"
```JSON
{
	"followid" : "might be empty",
	"userid" : "might be empty",
	"follower_userid" :"not empty"
}
```

- /user/post создает пост"POST"
- /user/post апдейтит пост"UPDATE"
- /user/post удаляет пост"DELETE"


- /user/comment создает комент"POST"
- /user/comment апдейтит комент"UPDATE"
- /user/comment удаляет комент"DELETE"
```JSON
{
	"commentid":"omitempty",
	"userid":"omitempty",
	"postid": "ОБЯЗ",
	"text": "ОБЯЗ",
}

```
- user/like принимает и вощвращает "POST"
- user/dislike только принимает"DELETE" 
```JSON
type Like struct {
	"likeid":"omitempty",
	"postid":"обяз",
	"userid":"omitempty",
}
```

- ДОБАВИТЬ МАССИВ JSON В routestat И СТАРТОВУЮ СТРАНИЦУ
```plantuml
@startuml
Users --> Followers : [1; *]
Posts --> Routes
Points --> Routes
RouteStats --> Routes
Comments --> Posts
Routes --> Users
Likes --> Posts
Likes --> Users

Users: userid,omitempty
Users: username,omitempty
Users: firstname
Users: lastname,omitempty
Users: email
Users: userpassword,omitempty
Users: phonenumber,omitempty
Users: encryptedpassword

Followers: id 
Followers: userid (тот на кого)
Followers: follower_userid (тот кто)

Comments: id bigint
Comments: userid int
Comments: postid bigint
Comments: commentary text

Posts: postid bigint
Posts: routeid bigint
Posts: date
Posts: photopath1
Posts: photopath2
Posts: photopath3
Posts: text

Likes: likeid
Likes: postid
Likes: userid

Routes: routeid
Routes: userid
Routes: date

Points: id 
Points: routeid
Points: latitude
Points: longitude
Points: point_index

RouteStats: id bigint
RouteStats: routeid
RouteStats: userid
RouteStats: workouttime
RouteStats: date
```