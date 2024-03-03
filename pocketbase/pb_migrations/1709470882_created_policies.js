/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "xzzzptkywwwyxtr",
    "created": "2024-03-03 13:01:22.766Z",
    "updated": "2024-03-03 13:01:22.766Z",
    "name": "policies",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "3ocuog7v",
        "name": "ingressAnnotations",
        "type": "select",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "asdf"
          ]
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
  const collection = dao.findCollectionByNameOrId("xzzzptkywwwyxtr");

  return dao.deleteCollection(collection);
})
