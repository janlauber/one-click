/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "a3bwaaqkbj0nfdu",
    "created": "2023-12-05 15:11:06.505Z",
    "updated": "2023-12-05 15:11:06.505Z",
    "name": "manifests",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "qze3nz4j",
        "name": "json",
        "type": "json",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
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
  const collection = dao.findCollectionByNameOrId("a3bwaaqkbj0nfdu");

  return dao.deleteCollection(collection);
})
