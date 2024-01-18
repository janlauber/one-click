/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("6jgxtwxj1fdutrh");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "6jgxtwxj1fdutrh",
    "created": "2023-11-13 19:24:27.223Z",
    "updated": "2023-11-13 21:07:54.859Z",
    "name": "tags",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "jc4mgten",
        "name": "name",
        "type": "text",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": ""
        }
      },
      {
        "system": false,
        "id": "ltusoej6",
        "name": "user",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "_pb_users_auth_",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "qdwfyqcb",
        "name": "color",
        "type": "select",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "default",
            "dark",
            "red",
            "green",
            "yellow",
            "indigo",
            "purple",
            "pink"
          ]
        }
      }
    ],
    "indexes": [],
    "listRule": "@request.auth.id != \"\" && @request.auth.id = user.id",
    "viewRule": "@request.auth.id != \"\" && @request.auth.id = user.id",
    "createRule": "@request.auth.id != \"\" && @request.auth.id = user.id",
    "updateRule": "@request.auth.id != \"\" && @request.auth.id = user.id",
    "deleteRule": "@request.auth.id != \"\" && @request.auth.id = user.id",
    "options": {}
  });

  return Dao(db).saveCollection(collection);
})
