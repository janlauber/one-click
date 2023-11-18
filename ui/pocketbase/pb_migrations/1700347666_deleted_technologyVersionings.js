/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("jcqj2nfdwwfnw3h");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "jcqj2nfdwwfnw3h",
    "created": "2023-11-13 19:31:07.040Z",
    "updated": "2023-11-13 19:31:07.040Z",
    "name": "technologyVersionings",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "40r2o62t",
        "name": "technology",
        "type": "relation",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "nwgspeyl1n5v30f",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "qclgu9iy",
        "name": "version",
        "type": "text",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": ""
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
