/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "22k6ts6gvnp46mc",
    "created": "2023-11-13 19:21:18.321Z",
    "updated": "2023-11-13 19:21:18.321Z",
    "name": "deployments",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "oo3v9faj",
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
        "id": "h2y9ouvn",
        "name": "technology",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "nwgspeyl1n5v30f",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
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
  const collection = dao.findCollectionByNameOrId("22k6ts6gvnp46mc");

  return dao.deleteCollection(collection);
})
