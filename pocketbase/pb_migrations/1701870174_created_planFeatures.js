/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "m9zfrl1ndxsejxd",
    "created": "2023-12-06 13:42:54.573Z",
    "updated": "2023-12-06 13:42:54.573Z",
    "name": "planFeatures",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "qcxoaiao",
        "name": "name",
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
    "listRule": "",
    "viewRule": "",
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("m9zfrl1ndxsejxd");

  return dao.deleteCollection(collection);
})
