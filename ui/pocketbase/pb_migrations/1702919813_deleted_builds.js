/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("vltwmhd9vavgi9p");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "vltwmhd9vavgi9p",
    "created": "2023-12-07 14:39:34.401Z",
    "updated": "2023-12-07 14:46:06.062Z",
    "name": "builds",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "zlzfmh8z",
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
        "id": "t2cfmff4",
        "name": "repository",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "w6ia8paw6verog6",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "xtafzhh6",
        "name": "project",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "7kff2zw80a7rmbu",
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
