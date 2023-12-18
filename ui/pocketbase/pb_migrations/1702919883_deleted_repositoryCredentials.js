/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("btx74c43k7bgm5a");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "btx74c43k7bgm5a",
    "created": "2023-12-07 14:37:09.382Z",
    "updated": "2023-12-07 14:38:04.198Z",
    "name": "repositoryCredentials",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "vjvc61pv",
        "name": "username",
        "type": "text",
        "required": false,
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
        "id": "uewqvme4",
        "name": "password",
        "type": "text",
        "required": false,
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
        "id": "bpca3yzz",
        "name": "user",
        "type": "relation",
        "required": false,
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
