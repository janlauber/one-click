/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("w6ia8paw6verog6");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "w6ia8paw6verog6",
    "created": "2023-12-07 14:37:30.916Z",
    "updated": "2023-12-07 14:37:56.850Z",
    "name": "repositories",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "2c0oe7p5",
        "name": "url",
        "type": "url",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "exceptDomains": [],
          "onlyDomains": []
        }
      },
      {
        "system": false,
        "id": "azzojjet",
        "name": "private",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "j6z0n9iz",
        "name": "repositoryCredential",
        "type": "relation",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "btx74c43k7bgm5a",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "fvz7rjkw",
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
