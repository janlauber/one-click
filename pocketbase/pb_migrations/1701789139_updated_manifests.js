/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("a3bwaaqkbj0nfdu")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "zrqmwzee",
    "name": "rollout",
    "type": "relation",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "22k6ts6gvnp46mc",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("a3bwaaqkbj0nfdu")

  // remove
  collection.schema.removeField("zrqmwzee")

  return dao.saveCollection(collection)
})
