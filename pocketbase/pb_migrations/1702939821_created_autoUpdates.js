/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "zpmsid5xlf5361s",
    "created": "2023-12-18 22:50:21.616Z",
    "updated": "2023-12-18 22:50:21.616Z",
    "name": "autoUpdates",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "yqgmgenf",
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
        "id": "vmrrtaqw",
        "name": "registry",
        "type": "select",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "github",
            "docker",
            "harbor"
          ]
        }
      },
      {
        "system": false,
        "id": "he7930nq",
        "name": "interval",
        "type": "text",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": "\\b(?:\\d+h)?(?:\\d+m)?(?:\\d+s)?\\b"
        }
      },
      {
        "system": false,
        "id": "weuvu4uj",
        "name": "regexPattern",
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
        "id": "pffoke7f",
        "name": "sortType",
        "type": "select",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "semver",
            "timestamp"
          ]
        }
      }
    ],
    "indexes": [],
    "listRule": null,
    "viewRule": null,
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s");

  return dao.deleteCollection(collection);
})
