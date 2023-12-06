/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("m9zfrl1ndxsejxd");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "m9zfrl1ndxsejxd",
    "created": "2023-12-06 13:42:54.573Z",
    "updated": "2023-12-06 13:50:06.131Z",
    "name": "planFeatures",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "rikp7sjv",
        "name": "options",
        "type": "json",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "joavjbxi",
        "name": "plan",
        "type": "relation",
        "required": true,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "s7nhljkrmzzu8y6",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      }
    ],
    "indexes": [],
    "listRule": "@request.auth.id != \"\"",
    "viewRule": "@request.auth.id != \"\"",
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
})
