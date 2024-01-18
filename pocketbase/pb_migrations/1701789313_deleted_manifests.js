/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("a3bwaaqkbj0nfdu");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "a3bwaaqkbj0nfdu",
    "created": "2023-12-05 15:11:06.505Z",
    "updated": "2023-12-05 15:14:34.816Z",
    "name": "manifests",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "zrqmwzee",
        "name": "rollout",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "22k6ts6gvnp46mc",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "qze3nz4j",
        "name": "json",
        "type": "json",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "nrwdaqpz",
        "name": "startDate",
        "type": "date",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "min": "",
          "max": ""
        }
      },
      {
        "system": false,
        "id": "grnqhwti",
        "name": "endDate",
        "type": "date",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "min": "",
          "max": ""
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
})
